README.md

$ minikube start
$ kubectl proxy		# optional if you want to see the dashboard
$ open http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/.
$ kubectl create serviceaccount cluster-admin-dashboard-sa
$ kubectl create clusterrolebinding cluster-admin-dashboard-sa \
  --clusterrole=cluster-admin \
  --serviceaccount=default:cluster-admin-dashboard-sa
$ kubectl get secret | grep cluster-admin-dashboard-sa
$ kubectl describe secret cluster-admin-dashboard-sa-token-tr4nd

### setup to use local image

$ export DOCKER_TLS_VERIFY="1" 
$ export DOCKER_HOST="tcp://192.168.49.2:2376"
$ export DOCKER_CERT_PATH="/Users/hermanwahyudi/.minikube/certs"
$ export MINIKUBE_ACTIVE_DOCKERD="minikube"

$ eval $(minikube docker-env)

$ docker build -t chatops:${tag} .

$ kubectl apply -f deployment-chatops.yml 

$ kubectl apply -f service-chatops.yml

$ minikube service go-chatops --url

$ ngrok http ${port}


### deploy configMap with imperative method

# first deploy
$ kubectl create configmap go-chatops-configmap --from-file config/config.json

# second and so on deployment
$ kubectl create configmap go-chatops-configmap --from-file config/config.json -o yaml --dry-run | kubectl replace -f -

# restart deployment
$ kubectl rollout restart deployment go-chatops