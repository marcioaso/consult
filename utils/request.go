package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Request faz uma requisição HTTP GET e retorna o corpo da resposta como um tipo genérico (interface{})
func Request(url string, headers map[string]string) (interface{}, error) {
	// Criar uma nova requisição GET
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Definir os cabeçalhos fornecidos
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Configurar o cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Fazer a requisição
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Ler o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Deserializar a resposta JSON
	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Retorna o corpo da resposta como um tipo genérico
	return result, nil
}
