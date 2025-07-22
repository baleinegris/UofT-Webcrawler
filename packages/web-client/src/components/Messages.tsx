import { MessageSchema } from "../schemas/Message";
import Message from "./Message";


export default function Messages({ messages }: { messages: MessageSchema[] }) {
  return (
    <div className="max-w-[80%] min-w-[80%] overflow-y-auto p-4 space-y-4">
        {messages.map((msg, index) => (
          <Message key={index} message={msg} />
        ))}
    </div>
  );
}