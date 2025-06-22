INSERT INTO features (name, type, geom, properties) VALUES
  (
    'Shibuya Station',
    'point',
    ST_GeomFromText('POINT(139.702 35.658)', 4326),
    '{"category": "station", "crowd": "high"}'
  ),
  (
    'Yamanote Line',
    'line',
    ST_GeomFromText('LINESTRING(139.702 35.658, 139.710 35.689)', 4326),
    '{"category": "rail", "color": "green"}'
  ),
  (
    'Imperial Palace',
    'polygon',
    ST_GeomFromText('POLYGON((139.752 35.685, 139.754 35.685, 139.754 35.688, 139.752 35.688, 139.752 35.685))', 4326),
    '{"category": "park", "area": "large"}'
  ),
  (
    'Sumida River Path',
    'line',
    ST_GeomFromText('LINESTRING(139.800 35.710, 139.810 35.700)', 4326),
    '{"category": "walking_path", "status": "open"}'
  ),
  (
    'National Stadium',
    'point',
    ST_GeomFromText('POINT(139.717 35.678)', 4326),
    '{"category": "stadium", "capacity": 68000}'
  );
