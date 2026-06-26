import { ContentPage } from "@/components/marketing/content-page";

export default function PortfolioPage() {
  return (
    <ContentPage
      title="Selected work that feels polished and credible"
      description="A premium portfolio layout helps convert interest into project inquiries."
      points={[
        { title: "Corporate website", description: "Minimal, strong typography, and clear value proposition." },
        { title: "Operations portal", description: "Dashboard-first interface for faster day-to-day usage." },
        { title: "Service marketplace", description: "Workflow-driven platform with clear onboarding." },
      ]}
    />
  );
}