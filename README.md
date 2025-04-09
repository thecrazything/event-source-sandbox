# Event sourcing sandbox

This projects is a sandbox for playing around with event-sourcing in a fullstack context using go, angular and kafka. Might add some Java too, who knows, I'm a wildcard.

## Architecture

This repo aims to implement a fullstack service without ever using request-reply or polling. Everything is reactive and event-driven.
As a starter, I've implemented a register and login flow (without a database at the time of writing, because 'can I save a hashed password' wasnt the point :) )

The basic structure for all requests from the browser is as follows:
1. Connect to websocket.
2. Websocket provides a SubscriberId for the connection
3. Any upstream content will close the websocket. It is downstream only. A re-connect means a new subscriberId.
4. Requests are done over HTTP, and are "rest-ish". All request must include the x-subscriber-id and x-request-id header. The requestId header is a client-provided Id for identifying requests that return on the socket.
5. ALL requets except the /session (we will get to why) responds with 202 Accepted (assuming its a correct request).
6. The response for the request then arrives over the websocket, with the requestId provided as an identifier.

Okay but.. why? That seems very complicated.

Because! Under the hood the above flow sent an event over kafka, that went through a series of services that eventually returned and triggered the websocket message.
At no point did any service in this flow need to have a wait timeout, or do any polling. If a service had gone down, as long as it comes up again the client would've gotten their message eventually.
Well, assuming their socket stayed connected. There are limitations of course. Might be ways around that, too!

Now this is all well and good for normal requests with some data, like the registration endpoint. But I mentiond that /session is special.
Well:

## Login flow

When we login to a webservice we expect it to be secure. In web, that means we need to follow some best practices to prevent our auth token from getting stolen.
This presents a problem for the above architecture because *websockets can't set cookies*. We want our response in a httponly, strict domain locked, secure cookie.

So for a login, the entire flow above applies with a few changes:
1. When doing the /login request, we get provided with a httponly strict cookie LOGIN-REQUEST-ID. This is an id that only lives for a minute. We will use this id to fetch our session later.
2. The events propagate as before, to the auth-service that checks our credentials and puts out a success or fail event.
3. Then we get back to the login gateway (not the websocket yet!). There, we use the loginRequestId to save our auth token into redis!
4. We then notify via the websocket to the client that "Hey, there is an auth cookie for you to get".
5. The user can now call /session which will read the LOGIN-REQUEST-ID cookie, use that to find the auth-token in redis (and then delete it)
6. We provide the auth-token as a cookie on the /session request! Login completed!

## Abstract the annoying parts in angular

Now you might think "oh no, when I do frontend dev I have to maintain a websocket and that is annoying". 

Luckily in angular we have rxjs and rxjs is magic.
In the service *request.service.ts* I have built a wrapper around HttpClient and the Websocket. This wrapper sends the request over http, then wait for the websocket before return the response.
So from the perspective of a client, it looks like a normal angular http request! You just to requestService.get<MyData>('/my-url').subscribe() and off you go! You wont even know the socket is there!
Your angular code is then also fully reactive! And of course for cases like a list where you want to be constantly subscribed you can handle the socket yourself and simpyl filter on your requestId.


## Disclaimer

Do not use this in production. Do not copy this into your own codebase without thinking critically about it. This is the first time I write go, I am sure these services do some really bad things.
I have no idea what happens if you run multiple intances of them. The kafa topics are also bad, they should probably have timesouts.

**This is intended as a proof-of-concept for the architecture.**
