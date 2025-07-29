# PoC
Note:
This repo is an Initial PoC code container. This repository will act as a playground for ideating the system for it's internal core members. 

Idealogy behind the system

AWS MGN Service is a major service used to perform lift and swift migrations of servers from various locations (on-prem or other clouds) to AWS. There are some other major lift and swift services out there in the industry too, each with their own limitation. We want to build a system that is cloud agnostic, modulur, open-source and light-weight. We could say the crux of this project is, "To build a system that can move servers from any destination to any cloud." 

We have decided to build the system to support migrating servers to AWS in our first phase, with a modular design pattern that would allow the system to accomadate future changes for replacing AWS with Azure or GCP, or any other cloud provider for that matter. We like to build the system with a "plug and play" design pattern, where the end-user would have a lot of freedom to switch lot of components to make the system work best for their use case. 

We have choosen "Golang" as the primary language to build the system with, as it allows the code to run in near native efficiency. We also choose Golang over other languages because of it's composition based design pattern rather than inheritance (i.e OOPS). "Change is the enemy of perfect design". The features supported by Golang, would be very helpful in near future should we want to incorporate a huge change. As mention above we want the system to support multi-cloud environment, which means we want to ensure the system will be able to accomodate any cloud-module we would build for it under an  abstraction layer.

Current Update: 
The "source_client" is an agent prototype that would run on the source-server side of things and the "mgn_Server" is an replication node prototype that mimics AWS MGN replication server. These current code is built to support linux systems. The code has been tested with ubuntu based machines. 

Agent: 
The agent in source would open a connection to the (TCP server hosted by) replication node and started sending data over the network. The data sent is the block-level data that is strored in the root EBS volume of the machine. This root volume is contains the OS and components that is required to boot the machine.  

Replication Server:
The replication servers hosts a TCP server that would grab data from source and write it to a EBS device. Once the data in root has been transfered from source, it can be used to spin up machines in the desired destination.

Create a Machine:
A new instance is created in the destination(AWS EC2) with an instanc-type that is similar to the source. This instance is stopped to remove it's root EBS drive. The EBS drive to which the replicated data is written to is attached as the root volume of this machine. If this instance is started then it would boot from the replicated EBS drive that contains all the data that was available in the source. 

Conclusion:
Now the current prototype involves lot of manuall interventions and hardcoded values in it. But this serves as the template we would need to work with. We aim to build the system that would be able to automate everything we did above.


List of items need to be completed:
1) Build an Agent module that will run on source with it's primary function for now being to open up connections to the replication node in the background, scan the device for it's specs, and replicate block level data from all mounted drives over the network to the replication node. Make sure the agent will accept instructions from the master and not act on it's own. The agent must only send data to the replication node specified to it, by the master, to avoid security risks of data being sent to somewhere else.

2) Build a replication node TCP server that would accept connections from the agent and store the data in EBS(volumeDevice). The replication node should be managed by the Master i.e the creation & destruction and it's configurations. It must only receive data from the agent source that it is tagged to. 

3) Build a master API server that would be able to manage the replication node and store lot of explicit details about the source. Store the data in a light-weight DB. Implement a mechanism that would allow this master to control the agent(start/stop the agent, communicate the location of replication node) and control the replication node(creation and mangement of these nodes). The master server is the most important component of the system, excercise caution while designing.

4) Build a Cloud server that will access AWS to provision and configure AWS resources. Appropriate permissions have to be give to the server for it to interact with AWS to manage resources. This is the part of the system we want to be modular and abstract so that we can swap it out for other cloud server that support some other cloud provider. (It's similar to having multiple cloud providers in terraform.)

5) Track changes in source and migrate those changes incrementally to the destination. (We think this is the most difficult part and the team needs to perform intense brainstroming to acheive this in the most effective way.)

6) Build front end, a client facing site that would allow ease of access to manage the system (front end can be vibe coded hehe (with some caution ofcourse) :-) )