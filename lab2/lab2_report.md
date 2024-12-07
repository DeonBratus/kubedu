# Лабораторная работа №2
**University**: [ITMO University](https://itmo.ru/ru/)\
**Faculty**: [FICT](https://fict.itmo.ru)\
**Course**: [Introduction to distributed technologies](https://github.com/itmo-ict-faculty/introduction-to-distributed-technologies)\
**Year**: 2024/2025\
**Group**: K4110c\
**Author**: Bratus Denis Alekseevich\
**Lab**: Lab2\
**Date of create**: 15.10.2024\
**Date of finished**: -\

---
## Описание лабораторной работы
Было создано два манифеста для деплоймента подов и для сервиса, соответственно [**deployment.yaml**](/lab2/deployment.yaml), [**service.yaml**](/lab2/service.yaml). 
В первом манифесте создается объект класса Deployment, устанвливаются метаданные и спецификации, наиболее интресными моментами является указание кол-ва реплик, в нашем случае будет создано две, а также установление стратегии k8s - RollingUpdate, что означет что k8s при изменениях будет обновлять постепенно по одному поду.


**deployment.yaml**
```yaml
apiVersion: apps/v1 
kind: Deployment # type of manifest is deployment, which can create and update pods
metadata:
  name: app-server # name of this deploymnet is app-server
spec: # specification of deployment
  replicas: 2 # creating 2 replicas
  selector: # using for connection between deployments and pods inside dep
    matchLabels: # set labels which help to find pods 
      app: app-server # name of match
  strategy: # set deployment srategy
    type: RollingUpdate # k8s will update only one pod at moment
    ...
```
Также в данном случае были созданы шаблоны для подов, внутри котороых они и описываются. Внутри шаблона указываются шаблонное имя, образ создаваемых контейнеров, и другие параметры. 
```yaml
...
template:  # templates for pods
    metadata:
      labels:
        app: app-server # template name of pods
    spec:
      containers: # info about containers
      - image: ifilyaninitmo/itdt-contained-frontend:master # image of conteiners
        imagePullPolicy: Always # k8s will try to launch latest version of img
        name: app-server # name of container in pods env
        ports:
        - containerPort: 3000 # open port 3000 
...
```
Далее описываются ресурсы выделямые для контейнеров минимальные и пределы
```yaml
...

resources: # set resources
    requests: # set minimums
        memory: "128Mi"
        cpu: "250m"
    limits: # set limits
        memory: "512Mi"
        cpu: "500m"
...
```
И в конце данного манифеста устанавливаются перменные окружения для выполнения поставленной задачи

```yaml
...
env: # set env values
- name: REACT_APP_USERNAME # set name for app
  value: "CHEL"
- name: REACT_APP_COMPANY_NAME # set company name for app
  value: "Mnogo_chelov"
```

Во втором документе описан сервис для созданного деплоймента, в нем устанавливается тип сервиса - **NodePort**, а также порты - порт внутри кластера **port** и внешний порт **targetPort**, открывать веб-приложение в дальнейшем необходимо будет именно через этот порт:


**service.yaml**
```yaml
apiVersion: v1 # set api version
kind: Service # set type Service
metadata:
  labels:
    app: app-server-service # set labels
  name: app-server-service # set name of service
spec: # specification
  type: NodePort # set type of service - NodePort
  ports: # set ports info
  - name: http # name of port
    port: 31000        # internal port inside cluster
    targetPort: 3000   # service will send trafic to this port
  selector:
    app: app-server
```

Далее необходимо применить оба манифеста, и прокинуть порты для сервиса:
```bash
minikube kubectl -- apply -f ./lab2/deployment.yaml   
minikube kubectl -- apply -f ./lab2/service.yaml
minikube kubectl -- port-forward service/app-server-service 3000:31000
```
После этого можно будет перейти по адресу [http://127.0.0.1:3000/](http://127.0.0.1:3000/)

![](/lab2/web-app-lab2.png)