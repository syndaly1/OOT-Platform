DELETE FROM contents
WHERE id::text LIKE '10000000-%'
   OR id::text LIKE '20000000-%'
   OR id::text LIKE '30000000-%'
   OR id::text LIKE '40000000-%'
   OR id::text LIKE '50000000-%';
