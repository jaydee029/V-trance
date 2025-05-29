# <img src="https://github.com/user-attachments/assets/b43ff8bc-ab93-4f09-9860-921b32f4bcf2" width="50px"> V-trance
V-Trance is a cost effective video transcoding and on demand streaming platform leveraging event driven architecture and microservices for scalable, asynchronous processing. It performs transcoding and/or stream creation, enabling multi format video output optimized for web delivery and adaptive bitrate streaming. 
<p align="center">
  <img src="https://github.com/user-attachments/assets/93f7ba84-b1d4-46e7-a923-43e8c61650f2" alt="First Image" style="width:46%; display:inline-block;"/>
  <img src="https://github.com/user-attachments/assets/27b4d0cf-f420-44eb-9453-43ba77b203a8" alt="Second Image" style="width:52%; display:inline-block;"/>
</p>

## Problem Statement & Introduction
Most existing services offered by cloud providers like AWS Elemental MediaConvert, Azure Media Services, or specialized platforms like Mux offer high-quality transcoding
and streaming infrastructure. However, they are also expensive, often based on per-minute pricing models that scale linearly with content volume. Additionally, traditional transcoding systems are monolithic and difficult to scale horizontally. They are often limited by single-node processing models, which become bottlenecks under high load.

The primary objective of this project is to design and build a scalable, modular, and costeffective platform for video transcoding and on-demand streaming that delivers
performance and flexibility comparable to commercial solutions, but at a significantly reduced cost.

## Features & Properties
- Provides various video manipulation services like transcoding, transrating as well as Adaptive bit rate video streaming using HLS(http live streaming) using FFmpeg.
- Allows user signp/login using stateless Access and Refresh tokens
- Microservices architecture implimented to enable scalability and asynchronous processing,
  - User service provided endpoints for user authorization and authentication.
  - trance-api service provides api endpoints for transcoding and stream creation
  - worker service does the heavy lifting , performing various transcoding and stream creation tasks
  - RabbitMQ used for asynchrnous job processing therby implementing event driven architecture
- PostgresSQL used as primary database. Traefik api gateway implemented for ease of access.
- Blackblaze B2 used for object storage, and Cloudflare CDN used for swift content delivery and reduced latency.

## Architecture and User flow
#### 1) High Level Design 
![image](https://github.com/user-attachments/assets/d77d4c70-d75a-4a33-a14a-4595b41f212e)

#### 2) Diagram describing user request flow
![image](https://github.com/user-attachments/assets/d4bdd358-04aa-4f8c-92d2-69f2f0738e5d)

#### 3) Diagram describing service request flow
![image](https://github.com/user-attachments/assets/3c663ebf-e3cd-48fd-b96f-c11c17a9401a)

  
