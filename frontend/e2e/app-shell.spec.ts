import { expect, test } from "@playwright/test";

test("home page shows the LeeTKOB app shell", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "LeeTKOB" })).toBeVisible();
  await expect(page.getByText("Project setup is in progress.")).toBeVisible();
});
