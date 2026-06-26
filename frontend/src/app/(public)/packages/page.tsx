import { ContentPage } from "@/components/marketing/content-page";

export default function PackagesPage() {
  return (
    <ContentPage
      title="Packages built for a premium buying journey"
      description="Simple, clear package framing helps buyers move from browsing to project submission faster."
      points={[
        { title: "Starter Package", description: "Great for landing pages and quick launches." },
        { title: "Growth Package", description: "Best for businesses that need multi-page delivery." },
        { title: "Pro Package", description: "Custom work for products and complex platforms." },
      ]}
    />
  );
}