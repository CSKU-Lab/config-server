// Ensure runner names stay unique inside configs database
const runnersDb = db.getSiblingDB('configs');

runnersDb.runCommand({
  createIndexes: 'runners',
  indexes: [
    {
      key: { name: 1 },
      name: 'unique_runner_name',
      unique: true,
    },
  ],
});
