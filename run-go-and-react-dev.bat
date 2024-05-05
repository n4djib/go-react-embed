@echo off
cd C:\_dev_projects\go-react-embed

echo Starting Go server...
start /B cmd /c "air"

echo Starting React server...
cd frontend
start /B cmd /c "npm run dev"
