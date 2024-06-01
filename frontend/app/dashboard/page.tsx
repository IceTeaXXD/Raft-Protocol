import React from 'react';
import Navbar from '../components/Navbar';

function Dashboard() {
    const terminalItems = [
        { id: 1, contents: ["$ Welcome to the terminal", "$ Clear"] },
        { id: 2, contents: ["$ TESTES", "$ make uff", "$ ls", "$ ls", "$ ls", "$ ls", "$ ls", "$ ls","$ ls", "$ ls", "$ ls", "$ ls", "$ ls", "$ ls", "$ ls"] },
        { id: 3, contents: ["$ TESTES", "$ make uff", "$ ls"] },
        { id: 4, contents: ["$ TESTES", "$ make uff", "$ ls"] },
    ];

    const statusItems = [
        { id: 1, contents: ["$ Status: Running"] },
        { id: 2, contents: ["$ Status: Running"] },
        { id: 3, contents: ["$ Status: Running"] },
        { id: 4, contents: ["$ Status: Running"] },
    ];

    return (
        <div>
            <Navbar />
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 m-2">
                {terminalItems.map((item, index) => (
                    <div
                        key={item.id}
                        className="relative bg-black text-green-500 p-4 font-mono h-64 overflow-auto"
                    >
                        <div className="absolute top-2 right-2 text-white text-sm z-10">
                            {statusItems[index].contents.map((status, idx) => (
                                <div key={idx}>{status}</div>
                            ))}
                        </div>
                        <div className="h-full overflow-auto pr-8">
                            {item.contents.map((content, idx) => (
                                <div key={idx}>{content}</div>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default Dashboard;
