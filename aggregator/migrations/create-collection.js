db.createCollection("news", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["title", "link", "published_at", "authors", "tags", "categories"],
      properties: {
        title: {
          bsonType: "string",
        },
        description: {
          bsonType: "string",
        },
        link: {
          bsonType: "string",
        },
        source: {
          bsonType: "string",
        },
        published_at: {
          bsonType: "date",
        },
        authors: {
          bsonType: "array",
          items: {
            bsonType: "string",
          },
        },
        tags: {
          bsonType: "array",
          items: {
            bsonType: "string",
          },
        },
        categories: {
          bsonType: "array",
          items: {
            bsonType: "string",
          },
        },
        content: {
          bsonType: "string",
        },
      },
    },
  },
});
