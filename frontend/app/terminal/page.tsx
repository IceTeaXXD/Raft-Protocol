'use client';

import Navbar from '../components/Navbar';
import { PromptInput } from '../components/PromptInput';
import { useState, useRef, useEffect } from 'react';

function TerminalPage() {
    const [prompts, setPrompts] = useState([{ input: '', output: '' }]);
    const inputRefs = useRef<HTMLInputElement[]>([]);

    const handleInput = (event: React.ChangeEvent<HTMLInputElement>, index: number) => {
        const newPrompts = [...prompts];
        newPrompts[index].input = event.target.value;
        setPrompts(newPrompts);
    };

    const processInput = (index: number) => {
        const processedInput = prompts[index].input;
        const newPrompts = [...prompts];
        newPrompts[index].output = processedInput;
        setPrompts([...newPrompts, { input: '', output: '' }]);
    };

    useEffect(() => {
        if (inputRefs.current.length > 0) {
            inputRefs.current[inputRefs.current.length - 1].focus();
        }
    }, [prompts]);
    return (
        <div>
            <Navbar />
            <main className="flex min-h-screen flex-col items-center justify-between p-24">
                <div className="w-[96%] h-[90%] md:w-[700px] md:h-[500px] lg:w-[800px] lg:h-[600px] flex flex-col border-8 border-black rounded-lg bg-black">
                    <div className="relative flex flex-col items-center">
                    </div>
                    <div className="w-full flex-1 bg-white pt-2 overflow-y-scroll scrollbar">
                        <h1 className="hidden md:block pl-2 my-2 font-terminal text-black">
                            Test
                        </h1>
                        <div>
                            {prompts.map((prompt, index) => (
                                <div key={index}>
                                    <PromptInput prompt={prompt.input} handleInput={handleInput} index={index} processInput={processInput} inputRefs={inputRefs} />
                                    {prompt.output && (
                                        <div className="pl-2 text-black font-terminal w-full break-words">
                                            Processed command: {prompt.output}
                                        </div>
                                    )}
                                </div>
                            ))}
                        </div>
                    </div>
                </div>
            </main>
        </div>
    )
}

export default TerminalPage