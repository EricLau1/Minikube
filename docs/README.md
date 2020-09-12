# Instalando o Kubernetes com o Minikube

O Minikube é uma ferramenta que facilita a execução local do Kubernetes. O Minikube executa um cluster Kubernetes de nó único dentro de uma Máquina Virtual (VM) no seu laptop para usuários que desejam experimentar o Kubernetes ou desenvolvê-lo diariamente.

## Minikube Features

O Minikube suporta os seguintes recursos do Kubernetes:

* DNS
* NodePorts
* ConfigMaps and Secrets
* Dashboards
* Container Runtime: Docker, CRI-O, and containerd
* Enabling CNI (Container Network Interface)
* Ingress

[Instalação](https://kubernetes.io/docs/tasks/tools/install-minikube/)

## Quickstart

Esta breve demonstração o guia sobre como iniciar, usar e excluir o Minikube localmente. Siga as etapas abaixo para iniciar e explorar o Minikube.

1. Inicie o Minikube e crie um cluster:

```bash
    minikube start
```

Para interarir com a interface no navegador execute um dos seguintes comandos:

```bash
    # Abre automaticamente o navegador padrão
    minikube dashboard

    # Imprime a url para ser usada no navegador
    minikube dashboard --url=true
```


> Se o minikube for instalado via `gcloud`, talvez seja necessário rodar o comando `gcloud components update` para atualizar o __kubectl__ e __minikube__.

A saída é semelhante a esta:

```bash
    Starting local Kubernetes cluster...
    Running pre-create checks...
    Creating machine...
    Starting local Kubernetes cluster...
```

Para obter mais informações sobre como iniciar seu cluster em uma versão específica do Kubernetes, VM ou tempo de execução do contêiner, consulte [Iniciando um cluster](https://kubernetes.io/docs/setup/learning-environment/minikube/#starting-a-cluster).

2. Agora, você pode interagir com seu cluster usando o kubectl. Para obter mais informações, consulte [Interagindo com seu cluster](https://kubernetes.io/docs/setup/learning-environment/minikube/#interacting-with-your-cluster).

Vamos criar uma implantação do Kubernetes usando uma imagem existente chamada echoserver, que é um servidor `HTTP` simples e expô-la na porta __8080__ usando `--port`.

```bash
    kubectl create deployment hello-minikube --image=k8s.gcr.io/echoserver:1.10
```

A saída é semelhante a esta:

```bash
    deployment.apps/hello-minikube created
```

3. Para acessar a implantação hello-minikube, exponha-a como um serviço:

```bash
    kubectl expose deployment hello-minikube --type=NodePort --port=8080
```

A opção `--type=NodePort` especifica o tipo do serviço.

A saída é semelhante a esta:

```bash
    service/hello-minikube exposed
```

4. O `hello-minikube` Pod foi lançado agora, mas você deve aguardar até que o Pod seja ativado antes de acessá-lo pelo Serviço exposto.

Verifique se o Pod está em funcionamento:

```bash
    kubectl get pod
```

Se a saída mostrar o `STATUS` como `ContainerCreating`, o Pod ainda estará sendo criado:

```bash
    NAME                              READY     STATUS              RESTARTS   AGE
    hello-minikube-3383150820-vctvh   0/1       ContainerCreating   0          3s
```

Se a saída mostrar o `STATUS` como `Running`, o Pod estará agora em funcionamento:

```bash
NAME                              READY     STATUS    RESTARTS   AGE
hello-minikube-3383150820-vctvh   1/1       Running   0          13s
```

5. Obtenha o URL do Serviço exposto para visualizar os detalhes do Serviço:


```bash
    minikube service hello-minikube --url
```

6. Para visualizar os detalhes do cluster local, copie e cole o URL obtido como saída no navegador.

A saída é semelhante a esta:

```bash
    Hostname: hello-minikube-7c77b68cff-8wdzq

    Pod Information:
    -no pod information available-

    Server values:
    server_version=nginx: 1.13.3 - lua: 10008

    Request Information:
    client_address=172.17.0.1
    method=GET
    real path=/
    query=
    request_version=1.1
    request_scheme=http
    request_uri=http://192.168.99.100:8080/

    Request Headers:
        accept=*/*
        host=192.168.99.100:30674
        user-agent=curl/7.47.0

    Request Body:
        -no body in request-
```

Se você não deseja mais que o Serviço e o cluster sejam executados, é possível excluí-los.

7. Exclua o serviço `hello-minikube`:

```bash
    kubectl delete services hello-minikube
```

A saída é semelhante a esta:

```bash
    service "hello-minikube" deleted
```

8. Exclua a implantação `hello-minikube`:

```bash
    kubectl delete deployment hello-minikube
```

A saída é semelhante a esta:

```bash
    deployment.extensions "hello-minikube" deleted
```

9. Pare o cluster local do Minikube:

```bash
    minikube stop
```

A saída é semelhante a esta:

```bash
    Stopping "minikube"...
    "minikube" stopped.
```

Para mais informações, veja [Parando um Cluster](https://kubernetes.io/docs/setup/learning-environment/minikube/#stopping-a-cluster):

Exclua o cluster local do Minikube:

```bash
    minikube delete
```

Para mais informações, veja [Excluindo um Cluster](https://kubernetes.io/docs/setup/learning-environment/minikube/#deleting-a-cluster)

## Gerenciando o seu Cluster


### Iniciando um Cluster

O comando `minikube start` pode ser usado para iniciar seu cluster. Este comando cria e configura uma Máquina Virtual que executa um cluster Kubernetes de nó único. Este comando também configura sua instalação do [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) para se comunicar com este cluster.

> Nota: Se você estiver atrás de um proxy da web, precisará passar essas informações para o comando `minikube start`: `shell https_proxy=<my proxy> minikube start --docker-env http_proxy=<my proxy> --docker-env https_proxy=<my proxy> --docker-env no_proxy=192.168.99.0/24` Infelizmente, definir apenas as variáveis de ambiente não funciona. O Minikube também cria um contexto de "minikube" e o define como padrão no kubectl. Para voltar a esse contexto, execute este comando: `kubectl config use-context minikube`.

### Especificando a versão do Kubernetes

Você pode especificar a versão do Kubernetes para Minikube a ser usada adicionando a cadeia de caracteres `--kubernetes-version` ao comando `minikube start`. Por exemplo, para executar a versão v1.18.0, você executaria o seguinte:

```bash
    minikube start --kubernetes-version v1.18.0
```

### Consultar erros:

```bash
    minikube logs
```