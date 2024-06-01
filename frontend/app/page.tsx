'use client';

import Navbar from './components/Navbar';

const Home = () => {
  return (
    <div>
      <Navbar />
      <div className='flex flex-col justify-center items-center'>
        <h1 className='text-3xl font-bold p-4'>Raft Visualization</h1>
        <p className='p-2'>by : SPG</p>
        <img src="https://media.licdn.com/dms/image/D5603AQEE9HQ1cd6rYQ/profile-displayphoto-shrink_200_200/0/1686828145381?e=2147483647&v=beta&t=vm9jl9KbmseusYXkAuCh8bau0RoQicVKUbMP5rwLKTc" alt="Description" height="800px" className="animate-pulse" />
      </div>
    </div>
  );
};

export default Home;
