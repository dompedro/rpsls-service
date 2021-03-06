CREATE (rock:Choice { name: "Rock" }),
(paper:Choice { name: "Paper" }),
(scissors:Choice { name: "Scissors" }),
(lizard:Choice { name: "Lizard" }),
(spock:Choice { name: "Spock" }),
(rock)-[:BEATS {with: "crushes"}]->(scissors),
(rock)-[:BEATS {with: "crushes"}]->(lizard),
(paper)-[:BEATS {with: "covers"}]->(rock),
(paper)-[:BEATS {with: "disproves"}]->(spock),
(scissors)-[:BEATS {with: "cuts"}]->(paper),
(scissors)-[:BEATS {with: "decapitates"}]->(lizard),
(lizard)-[:BEATS {with: "eats"}]->(paper),
(lizard)-[:BEATS {with: "poisons"}]->(spock),
(spock)-[:BEATS {with: "smashes"}]->(scissors),
(spock)-[:BEATS {with: "vaporizes"}]->(rock);
