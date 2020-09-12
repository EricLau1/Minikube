# Criando Chaves SSH  

## Comandos

- Criar uma chave:

```bash 
    # cria uma chave com criptografia rsa
    ssh-keygen

    # criando chave com criptografia diferente
    ssh-keygen -t dsa
```

Por padrão a chave gerada vai ficar no diretório `$HOME/.ssh/` no arquivo `id_rsa`.
Caso já exista a chave o programa perguntará se deseja sobreescrever a chave existente.
Senão quiser sobreescrever, escolha um diretório diferente para salvar a chave ou 
apenas coloque um nome diferente no arquivo em que a chave será salva.

O programa irá gerar dois arquivos com nomes iguais porém o arquivo com a chave privada 
não tem extensão, e o da chave pública possui a extensão `.pub`.

## Como Acessar um Servidor Remoto Com a Chave Pública

Para copiar a chave publica para um servidor via `ssh` será necessário saber a senha do servidor.

Execute o comando

```bash
    # Copia o arquivo para o direório home do usuario no servidor remoto
    scp id_rsa.pub usuario@endereco_ip:

    # exemplo
    scp id_rsa.pub eric@192.168.50.11:
```

Para que a chave funcione será necessário entrar no servidor e copiar o conteúdo do arquivo `id_rsa.pub`, 
para o arquivo `~/.ssh/authorized_keys`:

```bash
    cat id_rsa.pub >> ~/.ssh/authorized_keys
```

É possível copiar a chave para o servidor sem precisar entrar, com o comando:

```bash
    ssh-copy-id id_rsa.pub usuario@endereco_ip
```

## Configurar Agente SSH

```bash
    ssh-agent $SHELL

    # carrega as chaves privadas no agente
    ssh-add
```

Após digitar o `passphrase` uma vez é possível se conectar com o servidor 
remoto sem precisar de senha até que a máquina cliente reinicie ou o usuário
faça logout.
