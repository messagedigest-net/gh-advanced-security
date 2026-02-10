package services

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// InitConfig inicializa o Viper, lê o arquivo de config e variáveis de ambiente.
// Deve ser chamado pelo cmd/root.go no OnInitialize.
func InitConfig() {
	// 1. Encontrar diretório home
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 2. Configurar Viper para procurar .gh-advanced-security.yaml no $HOME
	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".gh-advanced-security")

	// 3. Ler Variáveis de Ambiente (prefixo GHAS_)
	// Ex: GHAS_PAGE=100 sobrescreve o config
	viper.SetEnvPrefix("GHAS")
	viper.AutomaticEnv()

	// 4. Definir Defaults (Valores padrão se nada mais for informado)
	viper.SetDefault("page", 20)
	viper.SetDefault("json", false)
	viper.SetDefault("debug", false)
	// viper.SetDefault("default_org", "minha-empresa") // Exemplo

	// 5. Tentar ler o arquivo de configuração
	// Ignoramos erro de "arquivo não encontrado" pois é opcional
	if err := viper.ReadInConfig(); err == nil {
		// Se quiser debugar: fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Helper para obter a organização padrão do arquivo de configuração
func GetDefaultOrg() string {
	return viper.GetString("default_org")
}
