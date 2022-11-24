package tado

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/gonzolino/gotado/v2"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type Tado struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	TadoUsername string `toml:"tado_username"`
	TadoPassword string `toml:"tado_password"`

	Log telegraf.Logger `toml:"-"`
}

func (plugin *Tado) SampleConfig() string {
	return sampleConfig
}

func (plugin *Tado) Init() error {
	if plugin.ClientID == "" {
		return fmt.Errorf("Client ID cannot be empty")
	}
	if plugin.ClientSecret == "" {
		return fmt.Errorf("Client Secret cannot be empty")
	}
	if plugin.TadoUsername == "" {
		return fmt.Errorf("Tado Username cannot be empty")
	}
	if plugin.TadoPassword == "" {
		return fmt.Errorf("Tado Password cannot be empty")
	}

	return nil
}

func (plugin *Tado) Description() string {
	return "Gather Tado readings"
}

func (plugin *Tado) Gather(a telegraf.Accumulator) error {
	ctx := context.Background()
	tado := gotado.New(plugin.ClientID, plugin.ClientSecret)

	user, err := tado.Me(ctx, plugin.TadoUsername, plugin.TadoPassword)
	if err != nil {
		return fmt.Errorf("tado: Unable to connect as '%s': %w", plugin.TadoUsername, err)
	}

	for _, h := range user.Homes {
		plugin.dumpHome(ctx, user, h, a)
	}
	return nil
}

func (plugin *Tado) dumpHome(ctx context.Context, u *gotado.User, h gotado.UserHome, a telegraf.Accumulator) {
	home, err := u.GetHome(ctx, h.Name)
	if err != nil {
		a.AddError(fmt.Errorf("Failed to find home '%s': %w\n", h.Name, err))
		return
	}

	zones, err := home.GetZones(ctx)
	if err != nil {
		a.AddError(fmt.Errorf("Failed to list zones in '%s': %w\n", h.Name, err))
		return
	}

	for _, z := range zones {
		plugin.dumpZone(ctx, u, home, z, a)
	}
}

func (plugin *Tado) dumpZone(ctx context.Context, u *gotado.User, h *gotado.Home, z *gotado.Zone, a telegraf.Accumulator) {
	state, err := z.GetState(ctx)
	if err != nil {
		a.AddError(fmt.Errorf("Failed to get zone state for '%s': %w\n", z.Name, err))
		return
	}

	tags := make(map[string]string)
	tags["home"] = h.Name
	tags["zone"] = z.Name
	fields := make(map[string]interface{})
	fields["setting"] = state.Setting.Temperature.Celsius
	fields["temperature"] = state.SensorDataPoints.InsideTemperature.Celsius
	fields["humidity"] = state.SensorDataPoints.Humidity.Percentage
	a.AddCounter("tado", fields, tags)
}

func init() {
	inputs.Add("tado", func() telegraf.Input {
		return &Tado{}
	})
}
