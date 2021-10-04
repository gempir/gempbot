# gempbot


- Frontend deployed on Vercel
- Most of backend on Vercel serverless functions
- bot docker container running on small scaleway server


## Setup

1. Login
2. Press Sub once (green button) and wait a moment (5s), this will subscribe the necesary webhooks to your channels

### Setup Channel Point Rewards

 - Go to Rewards and enable the Rewards you want.
 - make `gempbot` an editor in 7tv and bttv
 - Add Emote Blocks under the Blocks page

### Setup Prediction Management

- Join the bot to your channel from the sidebar


## Prediction Usage

## Starting

`!prediction Will nymn win this game?;yes;no;3m` --> yes;no;3m

`!prediction Will he win`                        --> yes;no;1m

`!prediction Will he win;maybe`                  --> maybe;no;1m
 
## Outcomes

`!outcome 1` or blue|first

`!outcome 2` or red|pink|second
 
 
## Locking/Aborting

`!prediction lock` --> will prevent any further submissions, otherwise just let the timer run out

`!prediction cancel` --> abort, if you typo'd or something