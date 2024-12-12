# Sistema de Clima por CEP com OTEL e Zipkin

Este projeto consiste em dois serviços (A e B) trabalhando em conjunto para fornecer informações de temperatura a partir de um CEP, bem como a integração com OpenTelemetry (OTEL) e Zipkin para observabilidade e tracing distribuído.

- **Serviço A**: Recebe um CEP via POST (JSON), valida se é um CEP válido (8 dígitos). Se válido, encaminha a solicitação para o Serviço B.
- **Serviço B**: Recebe o CEP, obtém a cidade via ViaCEP, consulta a temperatura atual via WeatherAPI, converte para Fahrenheit e Kelvin, e retorna o resultado junto com o nome da cidade.

Ambos os serviços enviam spans de tracing para um Otel Collector, que então exporta para o Zipkin, permitindo a visualização dos traces e spans.

## Pré-requisitos

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Como Iniciar o Ambiente

1. Clone o repositório e entre no diretório do projeto:
   ```bash
   git clone https://github.com/JMENDES82/Go-Expert-Open-Telemetry.git
   cd Go-Expert-Open-Telemetry
2. Ajuste a variável `WEATHER_API_KEY` no arquivo `docker-compose.yml` com a sua chave da [WeatherAPI](https://www.weatherapi.com/).

3. Inicie a aplicação e as ferramentas de observabilidade:
   ```bash
   docker-compose up --build

4. Testando a aplicação:,
    ```bash
    curl -X POST \
        -H "Content-Type: application/json" \
        -d "{\"cep\":\"29102571\"}" \
        http://localhost:8081/input

5. Serviços:   
    - **Serviço A**: [http://localhost:8081](http://localhost:8081)
    - **Serviço B**: [http://localhost:8080](http://localhost:8080)

    - **Zipkin**: [http://localhost:9411](http://localhost:9411)  
    Acesse esta URL para visualizar os traces coletados.
