'use client'

import { fetchHealth } from "@/features/tile/data/repository";

const a = () => {
  fetchHealth();
}
export default function Home() {
  
  return <>{a}</>;
}
