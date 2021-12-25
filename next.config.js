module.exports = {
    async rewrites() {
        return [
            {
                source: '/api/eventsub',
                destination: 'https://gempbot-api.gempir.com/api/eventsub',
            },
        ]
    },
}