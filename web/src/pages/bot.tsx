import { Bot as BotPage } from "../components/Bot/Bot";
import { initializeStore } from "../service/initializeStore";

export default function Bot() {
  return <BotPage />;
}

export const getServerSideProps = initializeStore;
