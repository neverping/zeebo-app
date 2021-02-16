import os

from diagrams import Cluster, Diagram
from diagrams.aws.compute import EKS
from diagrams.aws.network import ELB, InternetGateway, VPCRouter
from diagrams.onprem.network import Internet

from diagrams.k8s.network import Ing, SVC
from diagrams.k8s.compute import Deploy, Pod

graph_attr = {
    "fontsize": "45",
}

with Diagram("Zeebo-app Internet Traffic overview", show=False, filename="zeebo-aws-network-traffic", outformat="png", direction="TB", graph_attr=graph_attr):

    # Grouping step
    with Cluster("VPC"):
        igw = InternetGateway("IGW")
        with Cluster("Public Subnet"):
            with Cluster("Availability Zone A"):
                alb_leg_a = ELB("elb-interface-1")
            with Cluster("Availability Zone B"):
                alb_leg_b = ELB("elb-interface-2")
            with Cluster("Availability Zone C"):
                alb_leg_c = ELB("elb-interface-3")

        vpc_router = VPCRouter("Internal router")

        with Cluster("Private Subnet"):
            with Cluster("Availability Zone A"):
                eks1 = EKS("Kubernetes node 1")
            with Cluster("Availability Zone B"):
                eks2 = EKS("Kubernetes node 2")
            with Cluster("Availability Zone C"):
                eks3 = EKS("Kubernetes node 3")

    # Drawing step.
    Internet("Internet") >> igw

    igw >> alb_leg_a
    igw >> alb_leg_b
    igw >> alb_leg_c

    alb_leg_a >> vpc_router
    alb_leg_b >> vpc_router
    alb_leg_c >> vpc_router

    vpc_router >> eks1
    vpc_router >> eks2
    vpc_router >> eks3

with Diagram("Within K8s overview", show=True, filename="zeebo-inside-k8s-cluster", outformat="png", direction="TB", graph_attr=graph_attr):

    # Grouping step
    with Cluster("K8s"):
       ingress = Ing("zeebo-app.example.com - 80-443/tcp")

       svc_go = SVC("zeebo-go - 4458/tcp")
       svc_python = SVC("zeebo-python - 50051/tcp")

       deploy_go = Deploy("zeebo-go")
       deploy_python = Deploy("zeebo-python")

       pod_go_1 = Pod("zeebo-go-abc-node1")
       pod_go_2 = Pod("zeebo-go-def-node2")
       pod_go_3 = Pod("zeebo-go-ghi-node3")

       pod_python_1 = Pod("zeebo-py-jkl-node1")
       pod_python_2 = Pod("zeebo-py-mno-node2")
       pod_python_3 = Pod("zeebo-py-pqr-node3")

    # Drawing step.
    ELB("Internet") >> ingress >> svc_go >> deploy_go

    deploy_go >> pod_go_1
    deploy_go >> pod_go_2
    deploy_go >> pod_go_3

    pod_go_1 >> svc_python
    pod_go_2 >> svc_python
    pod_go_3 >> svc_python

    svc_python >> deploy_python

    deploy_python >> pod_python_1
    deploy_python >> pod_python_2
    deploy_python >> pod_python_3

