import { ContentPage } from "@/components/marketing/content-page";

export default function AboutPage() {
  return (
    <ContentPage
      title="Built by a team that values clarity and delivery"
      description="The platform is designed to feel serious, organized, and dependable for business buyers."
      points={[
        { title: "Design systems", description: "Visual consistency across all customer touchpoints." },
        { title: "Delivery discipline", description: "Structured progress from order to completion." },
        { title: "Trust and reliability", description: "A premium tone that reduces friction." },
      ]}
    />
  );
}