const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = {
    async rewrites() {
        return [
            {
                source: '/api_go/:path*',
                destination: 'http://localhost:8000/api/:path*', // Replace with your backend server URL
            },
        ];
    },
};
