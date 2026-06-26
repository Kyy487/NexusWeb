import { ContentPage } from "@/components/marketing/content-page";

export default function ContactPage() {
  return (
    <ContentPage
      title="Start a conversation"
      description="Use the platform to begin a project, ask a question, or request a custom scope."
      points={[
        { title: "Project inquiry", description: "Send requirements and move directly into the order workflow." },
        { title: "Support channel", description: "Follow the project and respond to updates from one place." },
        { title: "Payment guidance", description: "Clear invoice and payment steps reduce friction." },
      ]}
    />
  );
}