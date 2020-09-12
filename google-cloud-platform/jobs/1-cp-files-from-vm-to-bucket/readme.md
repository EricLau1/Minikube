# Copiando Arquivos de uma Máquina Virtual e Enviando Para o Bucket - Parte 1

Este exemplo considera que exista uma VM disponível com o Vagrant, e arquivos em pdf dentro dela.

- [Configurando VM para acessos via SSH/SCP](https://github.com/EricLau1/Vagrant-SSH-SCP-Examples)

Também será necessário ter uma conta do Google Cloud Platform configurada na máquina e o gcloud sdk para ter acesso aos comandos para interagir com o GCP.

- [Instalação do Gcloud SDK](https://cloud.google.com/sdk/install?hl=pt-br)

- Faça o login da conta pelo terminal:

```bash
    gcloud auth login

    # gera um arquivo com as credenciais no diretório: ~/.config/gcloud/application_default_credentials.json
    gcloud auth application-default login
```

- Configurar a conta e o projeto que serão usados:

```bash
    # mostrar as configurações atuais
    gcloud config list

    # setando a conta
    gcloud config set account <email@gmail.com>

    # setando o projeto
    gcloud config set project <project-id>
```

- Crie um bucket

```bash
    gsutil mb gs://my-bucket-2020  
```

A saída deve ser semelhante á esta:

```bash
Creating gs://my-bucket-2020/...
```

Referência: 

https://cloud.google.com/storage/docs/creating-buckets#storage-create-bucket-gsutil


- Faça um teste com este comando

```bash
    scp -o StrictHostKeyChecking=no -i ./ssh-keys/exemplo vagrant@192.168.50.11:/home/vagrant/*.pdf . &&  gsutil cp *.pdf gs://my-bucket-2020/ && rm *.pdf
```

> O comando copia o pdf dentro da vm via scp, para o diretório atual. Em seguida copia o pdf do diretório atual para o bucket. Após isso o pdf é removido do diretório.

Se o comando funcionou então é hora de configurar um Job para fazer isso sozinho!

Primeiro iremos construir uma imagem simples para fazer o primeiro teste.

- Crie um diretório chamado `Docker`:

```bash
    mkdir docker
```

- Crie um Dockerfile:

```Dockerfile
FROM alpine:3.7

RUN mkdir keys

COPY ssh-keys/exemplo /keys

RUN mkdir Uploads

RUN apk update && apk add openssh-client bash
```

- Crie o Job:

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: uploads
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: uploads
            image: uploads:1.0
            args:
            - /bin/sh
            - -c
            - scp -o StrictHostKeyChecking=no -i /keys/exemplo vagrant@192.168.50.11:/home/vagrant/*.pdf /Uploads && cd /Uploads && ls -la
          restartPolicy: OnFailure
```

### Executando o Job

```bash
    # iniciando o minikube
    minikube start

    eval $(minikube docker-env)

    # entrar na diretorio com Dockerfile
    cd docker

    # buildar a imagem
    sh build.sh

    # rodar imagem
    docker run --rm -it uploads:1.0 /bin/bash

    # Acessando a VM de dentro do Container
    ssh -o StrictHostKeyChecking=no -i keys/exemplo vagrant@192.168.50.11

    # Ativar o Job
    kubectl apply -f job.yml

    # Visualizar o job sendo executado
    minikube dashboard
```

O job está configurado para rodar a cada 1 minuto. 

É possível vizualizar os logs do job,
entrando no painel `Cron Jobs`, clicar no nome do job (`uploads`) e quando ele for executado, basta selecionar a opção __Logs__ na lateral direita.

A saída deve ser semelhante a esta:

```bash
Warning: Permanently added '192.168.50.11' (ECDSA) to the list of known hosts.
total 18256
drwxr-xr-x    1 root     root          4096 Sep 12 21:55 .
drwxr-xr-x    1 root     root          4096 Sep 12 21:55 ..
-rw-r--r--    1 root     root      18685041 Sep 12 21:55 progit.pdf
```

### Comandos uteis

```bash
    # visualizar cronjobs
    kubectl get cronjobs

    # Visualizar jobs
    kubectl get jobs

    # Remover cronjob
    kubectl delete cronjob uploads
```

Para que o job consiga enviar os arquivos para o bucket será necessário configurar uma conta de serviço no Google Cloud Platform. Para fazer isso basta acessar o painel `IAM e administrador`:

- Clicar na opção `Contas de serviço`.
- Clicar em `Criar Conta de Serviço`.
- selecionar a opção a role `Administrador do Storage`. 

Após criar a conta volte a tela inicial de `Contas de Serviço` e selecione a conta criada:

- Clicar em `Adicionar Chave`.
- Clicar em `Criar nova chave`.
- Selecionar `JSON`.
- Clicar em `Criar`.

Após baixar as credenciais no arquivo json, renomeie o arquivo para `credentials.json` e coloque-o num diretório seguro.

Antes de continuar destrua a imagem criada anteriormente:

```bash
  cd docker

  sh down.sh
```

Atualize o Dockerfile, adicionando as credenciais da conta de serviço e o sdk do gcloud:

```Dockerfile
FROM alpine:3.7

RUN mkdir keys

COPY ssh-keys/exemplo /keys
COPY credentials/credentials.json /keys

RUN mkdir Uploads

RUN apk update && apk add openssh-client python curl bash

RUN curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-228.0.0-linux-x86_64.tar.gz

RUN mkdir -p /usr/local/gcloud \
  && tar -C /usr/local/gcloud -xvf /google-cloud-sdk-228.0.0-linux-x86_64.tar.gz \
  && /usr/local/gcloud/google-cloud-sdk/install.sh

RUN rm /google-cloud-sdk-228.0.0-linux-x86_64.tar.gz

RUN curl -sSL https://sdk.cloud.google.com | bash && \
    echo 'source /root/google-cloud-sdk/path.bash.inc' >>~/.bashrc
```

Atualizando o __Job__ para enviar enviar os pdfs da máquina virtual para o bucket no Google Cloud Platform:

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: uploads
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: uploads
            image: uploads:1.0
            args:
            - /bin/sh
            - -c
            - scp -o StrictHostKeyChecking=no -i /keys/exemplo vagrant@192.168.50.11:/home/vagrant/*.pdf /Uploads && cd /Uploads && ls -la && /root/google-cloud-sdk/bin/gcloud auth activate-service-account --key-file=/keys/credentials.json && /root/google-cloud-sdk/bin/gsutil cp *.pdf gs://my-bucket-2020/
          restartPolicy: OnFailure
```

### Construindo a nova Imagem e rodando o Job

```bash
    # entrar na diretorio com Dockerfile
    cd docker

    # buildar a imagem
    sh build.sh

    # rodar imagem 
    docker run --rm -it uploads:1.0 /bin/bash

    # testar comando do gsutil
    gsutil --version

    # Ativar o Job
    kubectl apply -f job.yml

    # Visualizar o job sendo executado
    minikube dashboard
```

Basta repetir o mesmo precesso anterior para vizualizar os logs do job.

A saída deve ser semelhante a esta:

```bash
Warning: Permanently added '192.168.50.11' (ECDSA) to the list of known hosts.
total 18256
drwxr-xr-x    1 root     root          4096 Sep 12 22:44 .
drwxr-xr-x    1 root     root          4096 Sep 12 22:44 ..
-rw-r--r--    1 root     root      18685041 Sep 12 22:44 progit.pdf
Activated service account credentials for: [admin-bucket@curso-283912.iam.gserviceaccount.com]
Copying file://progit.pdf [Content-Type=application/pdf]...
/ [0 files][    0.0 B/ 17.8 MiB]                                                
/ [0 files][264.0 KiB/ 17.8 MiB]                                                
-
\
\ [0 files][ 11.6 MiB/ 17.8 MiB]                                                
|
| [1 files][ 17.8 MiB/ 17.8 MiB]                                                
/
Operation completed over 1 objects/17.8 MiB.  
```