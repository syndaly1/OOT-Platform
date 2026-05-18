import { allItems, collectionDataMap, getItemById } from "../data/catalog";

function delay(ms = 250) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

export async function getContentById(id) {
  await delay();
  return getItemById(id);
}

export async function getAllContent() {
  await delay();
  return allItems;
}

export async function getCollection(collection) {
  await delay();

  if (!collection || collection === "all") {
    return allItems;
  }

  return collectionDataMap[collection] || [];
}