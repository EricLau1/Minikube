## Copiando Arquivos de uma Máquina Virtual e Enviando Para o Bucket - Parte 2

### Esta parte irá criar um Job com Secrets 

> O valor de uma secret precisa ser em base 64.

Para transformar um valor qualquer em base 64 no linux, execute o comando:

```bash
    # Exemplos

    echo -n "chave privada" | base64 -w 0
    # saída
    Y2hhdmUgcHJpdmFkYQ==

    echo -n "credenciais" | base64 -w 0
    # saída
    Y3JlZGVuY2lhaXM=
```

Como ficará o arquivo de secrets:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: upload-keys
type: Opaque
data:
    pvtkey: Y2hhdmUgcHJpdmFkYQ==
    credentials: Y3JlZGVuY2lhaXM=
```

Como transformar o conteúdo de um arquivo para base 64:

```bash
# chave ssh privada
cat ssh-keys/exemplo | base64 -w 0

# credenciais do gcp
cat credentials/credentials.json | base64 -w 0
```

Após gerar o valor em base 64 basta setá-los para sua chave correspondente no arquivo de secrets.

Agora execute o comando para ativar a secret:

```bash
# Criar secret
kubectl apply -f secret.yml

# Visualizar secrets
kubectl get secrets

# Mostrar informações de um secret
kubectl describe secret upload-keys

# Remover secret
kubectl delete secret upload-keys
```

## Configurando a Secret dentro do Job

Alterar o arquivo `job.yml`:

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
          restartPolicy: OnFailure

          volumes:
          - name: uploads-volume
            secret:
              secretName: upload-keys
              defaultMode: 0400

          containers:
          - name: uploads
            image: uploads:1.0
            args:
            - /bin/sh
            - -c
            - scp -o StrictHostKeyChecking=no -i /keys/pvtkey vagrant@192.168.50.11:/home/vagrant/*.pdf /Uploads && cd /Uploads && ls -la && /root/google-cloud-sdk/bin/gcloud auth activate-service-account --key-file=/keys/credentials && /root/google-cloud-sdk/bin/gsutil cp *.pdf gs://my-bucket-2020/
            
            volumeMounts:
            - name: uploads-volume
              readOnly: true
              mountPath: /keys
```

Atualize o Dockerfile que agora não precisará mais ser criado com a chave privada e as credenciais do GCP:

```Dockerfile
FROM alpine:3.7

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

Contrua novamente a imagem:

```bash
  cd docker

  sh build.sh
```

Execute o Job:

```bash
kubectl apply -f job.yaml

minikube dashboard
```

Adicione mais alguns arquivos em pdf no máquina virtual para visualizar o resultado com vários arquivos.

A saída deve ser semelhante a esta:

```bash
Warning: Permanently added '192.168.50.11' (ECDSA) to the list of known hosts.
total 236344
drwxr-xr-x    1 root     root          4096 Sep 12 23:43 .
drwxr-xr-x    1 root     root          4096 Sep 12 23:43 ..
-rw-------    1 root     root       2954262 Sep 12 23:43 Computação Gráfica - Teoria e Prática - Eduardo Azevedo e Aura Conci.pdf
-rw-------    1 root     root     198557637 Sep 12 23:43 Computação Gráfica para Programadores Java.pdf
-rw-------    1 root     root       8800719 Sep 12 23:43 Mathematics For 3D Game Programming And Computer.pdf
-rw-r--r--    1 root     root       6252890 Sep 12 23:43 Programming-Language-Pragamatics.pdf
-rw-r--r--    1 root     root       2957681 Sep 12 23:43 distributed-services-with-go_B3.0.pdf
-rw-------    1 root     root       3787170 Sep 12 23:43 epdf.pub_mathematics-for-computer-graphics.pdf
-rw-r--r--    1 root     root      18685041 Sep 12 23:43 progit.pdf
Activated service account credentials for: [admin-bucket@curso-283912.iam.gserviceaccount.com]
Copying file://Computação Gráfica - Teoria e Prática - Eduardo Azevedo e Aura Conci.pdf [Content-Type=application/pdf]...
/ [0 files][    0.0 B/  2.8 MiB]                                                
/ [1 files][  2.8 MiB/  2.8 MiB]                                                
-
Copying file://Computação Gráfica para Programadores Java.pdf [Content-Type=application/pdf]...
- [1 files][  2.8 MiB/192.2 MiB]                                                
==> NOTE: You are uploading one or more large file(s), which would run
significantly faster if you enable parallel composite uploads. This
feature can be enabled by editing the
"parallel_composite_upload_threshold" value in your .boto
configuration file. However, note that if you do this large files will
be uploaded as `composite objects
<https://cloud.google.com/storage/docs/composite-objects>`_,which
means that any user who downloads such objects will need to have a
compiled crcmod installed (see "gsutil help crcmod"). This is because
without a compiled crcmod, computing checksums on composite objects is
so slow that gsutil disables downloads of composite objects.
- [1 files][  3.1 MiB/192.2 MiB]                                                
\
|
| [1 files][ 14.9 MiB/192.2 MiB]                                                
/
-
- [1 files][ 26.3 MiB/192.2 MiB]                                                
\
\ [1 files][ 37.6 MiB/192.2 MiB]                                                
|
/
/ [1 files][ 49.0 MiB/192.2 MiB]                                                
-
\
\ [1 files][ 60.3 MiB/192.2 MiB]                                                
|
| [1 files][ 71.7 MiB/192.2 MiB]                                                
/
-
- [1 files][ 83.0 MiB/192.2 MiB]                                                
\
|
| [1 files][ 94.3 MiB/192.2 MiB]   11.2 MiB/s                                   
/
/ [1 files][105.7 MiB/192.2 MiB]   11.2 MiB/s                                   
-
\
\ [1 files][117.0 MiB/192.2 MiB]   11.2 MiB/s                                   
|
| [1 files][128.4 MiB/192.2 MiB]   11.2 MiB/s                                   
/
-
- [1 files][139.7 MiB/192.2 MiB]   11.2 MiB/s                                   
\
|
| [1 files][151.1 MiB/192.2 MiB]   11.2 MiB/s                                   
/
/ [1 files][162.4 MiB/192.2 MiB]   11.2 MiB/s                                   
-
\
\ [1 files][173.8 MiB/192.2 MiB]   11.2 MiB/s                                   
|
/
/ [1 files][185.1 MiB/192.2 MiB]   11.2 MiB/s                                   
-
- [2 files][192.2 MiB/192.2 MiB]    8.4 MiB/s                                   
\
Copying file://Mathematics For 3D Game Programming And Computer.pdf [Content-Type=application/pdf]...
\ [2 files][192.2 MiB/200.6 MiB]    8.4 MiB/s                                   
|
| [2 files][197.9 MiB/200.6 MiB]    7.3 MiB/s                                   
/
/ [3 files][200.6 MiB/200.6 MiB]    4.3 MiB/s                                   
-
Copying file://Programming-Language-Pragamatics.pdf [Content-Type=application/pdf]...
- [3 files][200.6 MiB/206.5 MiB]    4.3 MiB/s                                   
\
\ [4 files][206.5 MiB/206.5 MiB]    3.1 MiB/s                                   
|
==> NOTE: You are performing a sequence of gsutil operations that may
run significantly faster if you instead use gsutil -m cp ... Please
see the -m section under "gsutil help options" for further information
about when gsutil -m can be advantageous.
Copying file://distributed-services-with-go_B3.0.pdf [Content-Type=application/pdf]...
| [4 files][206.5 MiB/209.4 MiB]    3.1 MiB/s                                   
| [5 files][209.4 MiB/209.4 MiB]    1.8 MiB/s                                   
/
Copying file://epdf.pub_mathematics-for-computer-graphics.pdf [Content-Type=application/pdf]...
/ [5 files][209.4 MiB/213.0 MiB]    1.8 MiB/s                                   
/ [6 files][213.0 MiB/213.0 MiB]    1.4 MiB/s                                   
-
Copying file://progit.pdf [Content-Type=application/pdf]...
- [6 files][213.0 MiB/230.8 MiB]    1.4 MiB/s                                   
\
\ [6 files][214.8 MiB/230.8 MiB]    1.8 MiB/s                                   
|
| [6 files][226.1 MiB/230.8 MiB]    4.1 MiB/s                                   
/
/ [7 files][230.8 MiB/230.8 MiB]    4.8 MiB/s                                   
-
Operation completed over 7 objects/230.8 MiB. 
```