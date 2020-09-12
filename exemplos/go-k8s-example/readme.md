# Golang com Kubernetes

1. Criar o `Dockerfile`:

```dockerfile
FROM golang

ENV GO111MODULE=on

WORKDIR go-k8s-example

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main main.go

EXPOSE 8080

CMD [ "./main" ]
```

2. Construir a imagem do __App__ no `Docker`:

```bash
    docker build -t go-k8s-example:1.0 .

    # verificar se imagem  roda sem erros:
    docker run -it --rm -p 8080:8080 go-k8s-example:1.0
```

> Para conseguir rodar esse imagem local no minikube será necessário construí-la após o `minikube start`.

Após deletar a imagem caso exista e, executar o `minikube start`, execute os comandos:

```bash
    eval $(minikube docker-env)

    # construa a imagem novamente
    docker build -t go-k8s-example:1.0 .
```


3. Criar o arquivo `pod.yaml` que irá definir um `Pod`:

```yml
apiVersion: v1
kind: Pod
metadata:
    name: go-k8s-example
spec:
    containers:
        - name: go-k8s-example
          image: go-k8s-example:1.0
          ports:
            - containerPort: 8080
```

4. Rodar o comando para criar o `Pod`:

```bash
    kubectl create -f pod.yaml 

    # COMANDOS PARA GERENCIAR UM POD
    
    # visualizar os pods
    kubectl get pod

    # visualizar informações de um pod
    kubectl describe pod go-k8s-example

    # deletar o pod
    kubectl delete pod go-k8s-example
```

5. Criar o arquivo `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-k8s-example
  labels:
    app: server
spec:
  replicas: 2
  selector:
    matchLabels:
      app: server
  template:
    metadata:
      labels:
        app: server
    spec:
      containers:
        - name: go-k8s-example
          image: go-k8s-example:1.0
          ports:
            - containerPort: 8080
```

> Arquivo de deployment define comportamentos para os pods. Neste caso definimos que o pod terá duas replicas/clones para que caso um deles seja destruído, o outro ocupe seu lugar e não derrube a aplicação.

6. Rodar o comando:

```bash
    kubectl create -f deployment.yaml

    # COMANDOS PARA GERENCIAR UM DEPLOYMENT

    # visualizar os deployment
    kubectl get deployment

    # visualizar informações de um deployment
    kubectl describe deployment go-k8s-example

    # deletar o deployment
    kubectl delete deployment go-k8s-example
```

Para expor o deployment execute o comando:

```bash
    kubectl expose deployment go-k8s-example --type=LoadBalancer --name=go-k8s-example

    # COMANDOS PARA GERENCIAR UM SERVICE

    # visualizar os deployment
    kubectl get service

    # visualizar informações de um deployment
    kubectl describe service go-k8s-example

    # deletar o deployment
    kubectl delete service go-k8s-example

```

> É possível utilizar o `--type=NodePort` no lugar de `--type=LoadBalancer`.

> Se o EXTERNAL-IP estiver como `<pending>` significa que não é possível obter um IP externo por que está rodando localmente. Normalmente esse IP externo e providenciado por algum serviço de Cloud como AWS ou Google. 

7. Obtendo a url do serviço:

```bash
    minikube service go-k8s-example --url
```

Cole a url no navegador para ver a aplicação funcionando.

Para visualizar como os pods se comportam quando ocorre uma falha no sistema faça um request para `/fail`.