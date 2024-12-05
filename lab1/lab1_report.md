# Лабораторная работа №1
**University**: [ITMO University](https://itmo.ru/ru/)\
**Faculty**: [FICT](https://fict.itmo.ru)\
**Course**: [Introduction to distributed technologies](https://github.com/itmo-ict-faculty/introduction-to-distributed-technologies)\
**Year**: 2024/2025\
**Group**: K4110c\
**Author**: Bratus Denis Alekseevich\
**Lab**: Lab1\
**Date of create**: 15.10.2023\
**Date of finished**: -\

---

## Описание лабораторной работы
Был создан манифест в файле [deployment.yaml](/lab1/deployment.yaml), где описан запускаемый "под". Ниже приведено содержимое данного манифеста:
```yaml
apiVersion: v1 # version of api k8s
kind: Pod # the type of object is Pod
metadata: 
  name: vault # the name of pod
  labels: 
    app: vault # the name of app, can be useful for searching
spec:  # this is specification for pod
  containers: # inforamtion about containers which will be launched in the pod
  - name: http # container name
    image: hashicorp/vault # image of container
    ports:
    - containerPort: 8200 # container listens to external port 8200  
```
После создание манифеста необходимо его применить, но перед этим необходимо установить и запустить *minikube*:\
**Установка minikube**
```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube && rm minikube-linux-amd64
```
**Запуск minikube**
```bash
minikube start
```
После его установки и запуска кластера необходимо применить манифест:
```bash
minikube kubectl -- apply -f deployment.yaml
```
После применения манифеста можно проверить запущен ли он и его статус:
```bash
minikube kubectl -- get pods
```
Вывести данная команда должна примерно следующее, где статус должен быть _running_, a "под" должен быть готов (1/1):
```bash
NAME    READY   STATUS    RESTARTS   AGE
vault   1/1     Running   0          123m
```
После этого необходимо создать сервис для созданного ранее пода, чтобы был доступ к контейнеру:
```bash
minikube kubectl -- expose pod vault --type=ClusterIP --port=8200
```
Также можно указать тип сервиса как **NodePort**, он тоже будет работать корректно.

Далее необходимо пробросить порт с помощью команды, и после чего можно будет перейти по адресу [_localhost:8200_](http://localhost:8200), где откроется нужна страница vault. 
```bash
minikube kubectl -- port-forward service/vault 8200:8200
```
Для получения токена доступа необходимо посмотреть логи запущеного пода, для этого существует следующая команда:
```bash
minikube kubectl -- logs vault
```
Далее нас интересует следующая строка:
```Root Token: hvs.*****MJ4A```
Послее ввода токена получится войти в систему, как на скриншоте снизу.
![](/lab1/vault_dashboard_screen.png)