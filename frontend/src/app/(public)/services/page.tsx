import { ContentPage } from "@/components/marketing/content-page";

export default function ServicesPage() {
  return (
    <ContentPage
      title="Digital Services for modern businesses"
      description="Service packages are structured for clear client expectations and premium presentation."
      points={[
        { title: "Web Development", description: "Fast, responsive, conversion-focused websites." },
        { title: "Application Interfaces", description: "Customer portals, dashboards, and internal tools." },
        { title: "Launch Support", description: "Project handoff, refinement, and go-live guidance." },
      ]}
    />
  );
}