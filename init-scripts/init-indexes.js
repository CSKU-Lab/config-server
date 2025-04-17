db = db.getSiblingDB("configs");

db.languages.createIndex({ id: 1 }, { unique: true });
db.compares.createIndex({ id: 1 }, { unique: true });
