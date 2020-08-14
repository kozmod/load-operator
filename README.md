   1 - Install metrics-server
   
    make metrics-server    
    
    note: metrics server has updatete args 
    
     containers:
            - name: metrics-server
              image: 'k8s.gcr.io/metrics-server/metrics-server:v0.3.7'
              args:
                - '--cert-dir=/tmp'
                - '--secure-port=4443'
                - '--metric-resolution=30s'
                - '--kubelet-preferred-address-types=InternalIP'
                - '--kubelet-insecure-tls'

   2 - run local minikube cluster
   
     make run ENABLE_WEBHOOKS=false