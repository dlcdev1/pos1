# Como executar o projeto

1. **Subir os containers dos serviços:**
   
   Execute o comando abaixo para iniciar todos os serviços necessários:
   
   ```sh
   docker-compose up -d
   ```
   
   Os seguintes serviços serão iniciados:
   - **serviço a**: responsável pelo input
   - **serviço b**: responsável pela orquestração
   - **zipkin**: responsável pelo trace
   - **otel-collector** : responsável por coletar os traces e enviar para o zipkin

2. **Realizar uma requisição POST:**
   
   Você pode utilizar o `curl` no Postman para enviar uma requisição:
   
   ```sh
   curl --location 'http://localhost:3000/cep' \
     --header 'Content-Type: application/json' \
     --data '{ "cep": "31748066" }'
   ```
   
   Ou, se preferir, utilize o site [ReqBin](https://reqbin.com/) com:
   - URL: `http://localhost:3000/cep`
   - Body:
     ```json
     {
       "cep": "31748066"
     }
     ```

3. **Visualizar o trace da requisição:**
   
   Acesse [http://localhost:9411/zipkin](http://localhost:9411/zipkin) para visualizar o trace das requisições.
   
   > **Nota:** As requisições só serão visualizadas após a segunda chamada, pois a primeira é utilizada para criar o trace.

