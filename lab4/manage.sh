if [ "$1" == "run" ]; then
    minikube kubectl -- apply -f configmap.yaml
    minikube kubectl -- apply -f replicaset.yaml
    minikube kubectl -- apply -f service.yaml
elif [ "$1" == "stop" ]; then
    minikube kubectl -- delete configmap --all
    minikube kubectl -- delete replicaset --all
    minikube kubectl -- delete service app-server-service 

fi