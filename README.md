# Notes

- Postman is used as frontend because I'm applying for the backend developer position
- There is no login system, instead users can request a token that is used to identify them - they can only see their own sessions
- As this is just a job test, there are some shortcuts and unfinished bits, but I tried to use the best possible practices.
- There are some example tests. Not full coverage, but enough for you get the idea
- Server config can be modified in app.properties
- The server returns basic data, which the client can use to calculate whatever it needs (eg duration = stop - start)
- sqlx is used instead of ORM systems, for example gorm, because I much prefer verbosity to obscurity
- K3s was chosen because of it's simple install method and ease of use
- Writing down every detail would take a lot of time, but I'd be glad to explain my reasons behind my decisions in the next meeting :)

# Requirements

- Go 1.16.7+
- K3s v1.21.3+
- Docker 20.10.8+

# Usage

Backend:
1. Install K3s and Docker, if you haven't already.
2. Run `$ ./build.sh <name of your image>`. This step can be skipped if you want to use my prebuilt image
3. Run `$ ./deploy.sh <name of your image>`. Alternatively, use `beranabus/pento-challenge` which is my prebuilt image.

Frontend:
1. Import postman collection
2. Before using any endpoints, call the `Get token` endpoint to "log in"
3. Start a session. Postman will store its uuid and use it in all next requests (stop, (re)name). Starting a session will automatically save it
4. You can now stop, or (re)name the session.
5. Listing usage:
   - `/list`: All sessions
   - `/list/day`: Last 24 hours
   - `/list/week`: Last 7 days
   - `/list/month`: Last 4 weeks

---
# Pento tech challenge

Thanks for taking the time to do our tech challenge.

The challenge is to build a small full stack web app, that can help a freelancer track their time.

It should satisfy these user stories:

- As a user, I want to be able to start a time tracking session
- As a user, I want to be able to stop a time tracking session
- As a user, I want to be able to name my time tracking session
- As a user, I want to be able to save my time tracking session when I am done with it
- As a user, I want an overview of my sessions for the day, week and month
- As a user, I want to be able to close my browser and shut down my computer and still have my sessions visible to me when I power it up again.

## Getting started

You can fork this repo and use the fork as a basis for your project. We don't have any requirements on what stack you use to solve the task, so there is nothing set up beforehand.

## Timing

- Don't spend more than a days work on this challenge. We're not looking for perfection, rather try to show us something special and have reasons for your decisions.
- Get back to us when you have a timeline for when you are done.

## Notes

- Please focus on code quality and showcasing your skills regarding the role you are applying to.
