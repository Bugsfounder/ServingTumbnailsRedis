"use client"
import React, { useEffect, useState } from "react"
import thumbImg from "./thumbnails/image.jpg"
import Image from 'next/image'
export default function Home() {
  const [thumbnail, setThumbnail] = useState(thumbImg)
  const ApiTest = async () => {
    const request = await fetch('/api');
    const response = await request.json();
    console.log(response);
  }
  useEffect(() => { ApiTest() }, [])
  return (
    <main>
      <h1 className="text-xl font-extrabold">Hello From react</h1>
      <Image
        src={thumbnail}
        alt="Picture of the author"
        width="350px"
        height="300px"
        placeholder="blur" // placeholder="empty" 
      />
    </main>
  )
}