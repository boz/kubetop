FROM busybox 

# kubectl run -it --image=kubetop --restart=Never --rm kubetop

ADD ./kubetop-linux /kubetop

ENTRYPOINT ./kubetop
