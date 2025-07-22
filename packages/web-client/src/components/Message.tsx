import { MessageSchema } from "../schemas/Message";

export default function Message( { message }: { message: MessageSchema } ) {
  const isUser = message.role.toLowerCase() === 'user';
  const isAgent = message.role.toLowerCase() === 'agent';
  
  return (
    <div className={`flex mb-4 ${isUser ? 'justify-end' : 'justify-start'}`}>
      <div className={`max-w-[80%] ${isUser ? 'bg-blue-500 text-white rounded-br-none' : 'bg-gray-200 text-black rounded-bl-none'} p-3 rounded-lg`}>
        <p className="text-sm leading-relaxed"><span className="font-bold">{message.role}: </span>{message.content}</p>
        {message.timestamp && (
          <span className="text-xs opacity-50 block mt-1">{message.timestamp}</span>
        )}
      </div>
    </div>
  );
}
