# Pub Sub

## Intro
- async messaging
- 4 components
    - Messages - message is communication data sent from sender to reciever
    - Topics - groups that messages belong to an intermediary channel b/n 
      sender and receivers
    - Subscribers - message recipient subscribes to topics
    - Publishers - publish messges to topics

- Push Delivery - async event notifications when messages are published
- Fanout - replicated to all subscribers
- Filtering - only let subscribers get messages they want
- Multiplexing - message streams are multiplexed in a consistent manner

- Advantages
    - Elminate polling
    - Dynamic targetting - discovery of services becomes easier
      dont have to keep track of all subscribers to send messages to
    - Decouple systems - allow for different subs and pubs to scale independently
    - Simplify communication - removes point to point communication and intergration
    - Durability - at least once delivery and durabiltiiy of messages by storing copies
    - Security - topics authenticate publishers and allow encrypted endpoints

- Messsage Queue vs pub/sub
    - Message queue - messages stored in the queue until consumed and deleted
        - require the sender to know who they are exchanging messages with.
        - order messages which may cause bottlenecks
    - pubsub allows more flexibility - same message can be processes by multiple 
      subscribers.


