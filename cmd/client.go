package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	iclient "github.com/marianoalvarez/freelow/internal/client"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Gestiona clientes",
}

var clientAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Registra un nuevo cliente",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		if name == "" {
			parts := strings.Split(id, "-")
			for i, p := range parts {
				if len(p) > 0 {
					parts[i] = strings.ToUpper(p[:1]) + p[1:]
				}
			}
			name = strings.Join(parts, " ")
		}
		if color == "" {
			color = "4"
		}
		if !iclient.ValidColor(color) {
			return fmt.Errorf("color inválido %q: usa un número ANSI (0-255) o hex (#RGB / #RRGGBB)", color)
		}

		cfg, err := iclient.Load()
		if err != nil {
			return err
		}
		for _, c := range cfg.Clients {
			if c.ID == id {
				return fmt.Errorf("client %q already exists", id)
			}
		}
		cfg.Clients = append(cfg.Clients, iclient.Client{
			ID:    id,
			Name:  name,
			Color: color,
		})
		if cfg.Active == "" {
			cfg.Active = id
		}
		if err := iclient.Save(cfg); err != nil {
			return err
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf("✓ Cliente %q añadido", id)))
		return nil
	},
}

var clientListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos los clientes",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := iclient.Load()
		if err != nil {
			return err
		}
		if len(cfg.Clients) == 0 {
			fmt.Println("No hay clientes. Usa: freelow client add <id>")
			return nil
		}
		activeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
		for _, c := range cfg.Clients {
			if c.ID == cfg.Active {
				fmt.Println(activeStyle.Render(fmt.Sprintf("▶ %s  %s  (activo)", c.ID, c.Name)))
			} else {
				fmt.Println(normalStyle.Render(fmt.Sprintf("  %s  %s", c.ID, c.Name)))
			}
		}
		return nil
	},
}

var clientSwitchCmd = &cobra.Command{
	Use:   "switch <id>",
	Short: "Cambia el cliente activo",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		cfg, err := iclient.Load()
		if err != nil {
			return err
		}
		if _, err := cfg.FindByID(id); err != nil {
			return err
		}
		cfg.Active = id
		if err := iclient.Save(cfg); err != nil {
			return err
		}
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		fmt.Println(style.Render(fmt.Sprintf("✓ Cliente activo: %s", id)))
		return nil
	},
}

func init() {
	clientCmd.AddCommand(clientAddCmd)
	clientCmd.AddCommand(clientListCmd)
	clientCmd.AddCommand(clientSwitchCmd)
	clientAddCmd.Flags().StringP("name", "n", "", "Nombre del cliente")
	clientAddCmd.Flags().StringP("color", "c", "4", "Color ANSI (0-255) o hex (#RGB / #RRGGBB)")
}
