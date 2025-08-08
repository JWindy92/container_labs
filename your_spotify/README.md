
Startup steps Ngrok:

1. Start containers: `docker compose up -d`
2. In a terminal, start ngrok tunnel on 443: `ngrok http 443`
3. Add ngrok endpoint + "oauth/spotify/callback" to Redirect URIs section in spotify application config. Ex: `https://0d5119e965ec.ngrok-free.app/api/oauth/spotify/callback`
4. Access app via ngrok endpoint


Startup steps Localtunnel:

1. Start containers: `docker compose up -d`
2. In a terminal, start ngrok tunnel on 443: `lt --port 8081 --subdomain yourspotify` (supports custom subdomain for free)
3. Add endpoint + "oauth/spotify/callback" to Redirect URIs section in spotify application config. Ex: `https://yourspotify.loca.lt/api/oauth/spotify/callback`
4. Access app via https://yourspotify.loca.lt

Spotify dashboard: https://developer.spotify.com/dashboard/525570b545c94b0f9173226823297728

Ngrok: https://dashboard.ngrok.com/

Localtunnel: https://localtunnel.github.io/www/

Your Spotify image docs: https://github.com/Yooooomi/your_spotify

