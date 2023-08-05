
"use client"
import React, { useState, useEffect } from 'react';

const Subscriber = () => {
    const [message, setMessage] = useState('');

    useEffect(() => {
        console.log('EventSource connecting...');
        const eventSource = new EventSource('/api/subscribe');

        eventSource.onmessage = (event) => {
            console.log('EventSource connected.');
            const message = event.data;
            console.log("message", message);
            setMessage(message);
        };

        return () => {
            eventSource.close();
        };
    }, []);

    return (
        <div>
            <h2>Redis Subscription Example</h2>
            <ul>
                <p>{message}</p>
                <img src={message} alt="" />
            </ul>
        </div>
    );
};

export default Subscriber;
