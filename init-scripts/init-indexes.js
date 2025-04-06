db = db.getSiblingDB("configs");

db.languages.createIndex({ id: 1 }, { unique: true });
