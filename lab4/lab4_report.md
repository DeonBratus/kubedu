# Лабораторная работа №4
**University**: [ITMO University](https://itmo.ru/ru/)\
**Faculty**: [FICT](https://fict.itmo.ru)\
**Course**: [Introduction to distributed technologies](https://github.com/itmo-ict-faculty/introduction-to-distributed-technologies)\
**Year**: 2024/2025\
**Group**: K4110c\
**Author**: Bratus Denis Alekseevich\
**Lab**: Lab4\
**Date of create**: 07.12.2024\
**Date of finished**: -\

---

## Описание лабораторной работы

В данной работе имеется 5 файлов [configmap.yaml](/lab4/configmap.yaml), [replicaset.yaml](/lab4/replicaset.yaml), [service.yaml](/lab4/service.yaml) и пулы ip-адресов для calico [ru-north-pool.yaml](/lab4/ru-north-pool.yaml), [ru-south-pool.yaml](/lab4/ru-south-pool.yaml).
Первые три файла не изменялись, только лишь в service.yaml была добавлена строчка ```nodePort: 3200``` для проброса порта, чтобы не писать команду ``` minikube kubectl -- forward-port...```.

---
### Запуск minikube с двумя нодами и calico в качестве CNI
---
Спрева был запущен minikube с необходимыми настройками, перед этим были очищены все прошлые профили и контейнеры minikube командой ``` minikube delete```, иначе *Calico* не заработает в создаваемых контейнерах. 

После чего необходимо было запустить minikube с указанием в качестве *CNI (Container Network Interface)* calico. 
```bash
minikube start --cni=calico
```
В документации от tigera указан также опция ```--network-plugin=cni```, однако при запуске minikube выводится предупреждение о том, что стоит использовать команду в том формате, что была указана выше.

После этого необходимо проверить работают ли поды calico, для этого необходимо ввести следующую команду:
```bash
watch minikube kubectl -- get pods -l k8s-app=calico-node -A
```
Результат должен быть примерно следующим:
![](/lab4/imgs/calico_pods.png)

После проверки calico необходимо запустит вторую ноду minikube, для этого можно воспользоваться командой
```bash
minikube node add
```
Для проверки можно ввести команду ``` minikube kubectl -- get nodes```. Увидеть должны следующее:
![](/lab4/imgs/minikube_nodes.png)
Итого был запущено две ноды и использова CNI calico.
### Работа c node labels и IPPools
Далее необходимо было установить метки для нод, в качестве меток были по признаку географического положения *ru-north* и *ru-south*. Для указания меток необходимо ввести следующие команды:
```bash
minikube kubectl -- node minikube zone=ru-north
```
```bash
minikube kubectl -- label node minikube-m02 zone=ru-south
```
Метки должны были примениться к нодам, для их проверки необходимо ввести команду
```bash
minikube kubectl -- get nodes --show-labels
```
В результате долнжы поулчить следующее:
![](/lab4/imgs/labels_of_nodes.png)
Далее необходимо было написать два IPPool-манифеста, **ru-north-pool.yaml** и **ru-south-pool.yaml**. Ключевыми отличями этих файлов являются поля *cidr* и *nodeSelector*. В cidr указываются те ip-адреса которые будут назначены подами из той или иной ноды, из северной и южной соответсвено. А в nodeSelector указывается ранее установленные labels.

**ru-north-pool.yaml** 
```yaml
apiVersion: crd.projectcalico.org/v1
kind: IPPool
metadata:
  name: north-zone-ippool
spec:
  cidr: 192.168.1.0/24
  ipipMode: Always
  natOutgoing: true
  nodeSelector: zone == "ru-north"
```
**ru-south-pool.yaml** 
```yaml
apiVersion: crd.projectcalico.org/v1
kind: IPPool
metadata:
  name: south-zone-ippool
spec: 
  cidr: 192.168.2.0/24
  ipipMode: Always
  natOutgoing: true
  nodeSelector: zone == 'ru-south'
```
Далее остается только применить эти два манифеста
```bash
minikube kubectl -- apply -f ./lab4/ru-north-pool.yaml
```
```bash
minikube kubectl -- apply -f ./lab4/ru-south-pool.yaml
```
Для проверки примененных IPPools можно использовать команду ``` minikube kubectl -- get ippools.crd.projectcalico.org  ```.

### Запуск веб-приложения и проверка работы

Для запуска веб-приложения с двумя репликами можно просто запустить скрипт ``` ./manage.sh run```, он примет необходимые манифесты, в результате чего появятся replicaset, c двумя репликами контейнера веб-приложения, один конфигмап с лежащими в нем данными пользователя и комании, а также сервис, к через который будет возможен доступ к приложению по адресу < IP-MINIKUBE >:32000, в моем случае это http://192.168.49.2:32000.

Выполнив команду ```minikube kubectl -- get pods -o wide``` можно узнать инофрмацию о подах.
![](/lab4/imgs/2cont_on_nodes_with_diff_ips.png)
Здесь интересуют имена и IP-адреса подов. Также здесь видно что поды запущены на разных нодах, это делает k8s автоматически балансировщиком нагрузки. 

После чего необзодимо перейти на страницу по ранее упомнянутому адресу, и несколько раз придется обновить страницу, чтобы увидеть, что имена и ip-адреса подов изменяются. На скриншоте ниже видны разные имена и адреса.
![](/lab4/imgs/really_2_cont's_name.png)
По моему мнению это происходит по той причине что *kubernetes* автоматически **переключает нагрзуку с одного на другой под**, по своим внутренним механизмам и когда пользователь делает запрос к сервису, то k8s **обращается к одному из подов**, а у них как было ранее выяснено разные имена и ip-адреса, связи с чем и видим их на странице.

Далее необходимо было проверить пингуется ли один под другим, для этго необходимо "войти" в под и пропинговать другой.

```minikube kubectl -- exec -it app-replicaset-jjnp5 -- /bin/sh```

И внутри конейнера был пропингован другой контейнер ```ping 192.168.1.66```. Результат работы представлен ниже на скриншоте:
![](/lab4/imgs/ping_cont2cont.png)