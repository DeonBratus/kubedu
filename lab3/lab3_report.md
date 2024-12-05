# Лабораторная работа №3
**University**: [ITMO University](https://itmo.ru/ru/)\
**Faculty**: [FICT](https://fict.itmo.ru)\
**Course**: [Introduction to distributed technologies](https://github.com/itmo-ict-faculty/introduction-to-distributed-technologies)\
**Year**: 2024/2025\
**Group**: K4110c\
**Author**: Bratus Denis Alekseevich\
**Lab**: Lab3\
**Date of create**: 05.12.2023\
**Date of finished**: -\

---
## Описание лабораторной работы
В данной работе 4 файла [replicaset.yaml](/lab3/replicaset.yaml), [service.yaml](/lab3/service.yaml), [configmap.yaml](/lab3/configmap.yaml), [ingress.yaml](/lab3/ingress.yaml).
В первом описываются создаваемые реплики подов, второй является сервисом для подов, он такой же как и во 2-й лабораторной работе, в configmap указаны переменные окружения, а в ingress описывается защищенное подключение к веб-приложению.
Стоит расмотреть сначала **configmap.yaml**, здесь в *data* указаны необходимые поля.:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-configmap
data:
  REACT_APP_USERNAME: "CHEL"
  REACT_APP_COMPANY_NAME: "Mnogo_chelov"
```
Далее необходимо было применить данный configmap:
```yaml
    minikube kubectl -- apply -f ./lab3/configmap.yaml
```
Необходимо включить аддон ingress: 
```bash
minikube addons enable ingress
```
После этого необходимо создать TSL-ключ и TSL-сертификат, сделать это можно с помощью следующей команды:
```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=myapp.local/O=myapp"
```
Будут получены *tls.crt* и *tls.key* файлы, их в дальнейшем необходимо будет использовать для создания секрета в minikube. Для этого необходимо выполнить следующую команду, указав пути к ключу и сертификату:
```bash
minikube kubectl -- create secret tls app-tls-secret --key tls.key --cert tls.crt
``` 
Для проверки был ли создан ключ можно использовать команду ``` minikube kubectl -- get secrets ```.

Далее необходимо создать файл **ingress.yaml**, в котором будут указаны параметры.
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: app-ingress
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true" 
spec:
  tls:
  - hosts:
    - app.local
    secretName: app-tls-secret
  rules:
  - host: app.local
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-server-service
            port:
              number: 3000
```
После этого необходимо лишь применить файл командой
```bash
minkube kubectl -- apply -f ./lab3/ingress.yaml
```
В заключении остается прописать домен в hosts, для этого сначала необходимо узнать IP-адрес minikube, это можно сделать с помощью команды ``` minikube ip``` и вставить его в файл /etc/hosts в формате "IP-MINIKUBE app.local".
Далее необходимо будет перейти по адресу https://app.local и проверить TSL-сертификацию. 

Скриншот работоспособности TSL-сертификации в веб-приложении
![](/lab3/screen_crt_is_here_lab3.png)

Описание сертификата
![](/lab3/more_info_about_crt_lab3.png) 