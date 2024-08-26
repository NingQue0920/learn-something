#!/bin/bash

# Function to display usage information
show_usage() {
    echo "Usage: $0 <repository_name>"
    echo "  repository_name: The path to the Git repository"
    echo
    echo "Example: $0 project"
}


# 函数：检查上一个命令是否成功
check_status() {
    if [ $? -ne 0 ]; then
        echo "Error: $1"
        exit 1
    fi
}

# 函数：处理本地修改冲突
handle_local_changes() {
    if git diff --quiet; then
        return 0
    fi
    log "检测到本地有未提交的修改。"
    log "丢弃本地修改以解决冲突..."
    git reset --hard
    git clean -fd
    log "本地修改已丢弃。"
}

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# 仓库名 - 仓库路径映射
declare -A repo_map
repo_map["alphamon"]="/data/source/alphamon"
repo_map["alphaservice"]="/data/source/alphaservice"
repo_map["alphadb"]="/data/source/alphadb"
repo_map["client"]="/data/source/client"
repo_map["alpha"]="/data/source/alpha"
repo_map["influxdb"]="/data/source/influxdb"




# 检查是否提供了仓库路径
if [ $# -eq 0 ]; then
    echo "Error: Repository path is required."
    show_usage
    exit 1
fi

REPO_NAME="$1"
REPO_PATH="${repo_map[$REPO_NAME]}"

# 检查提供的路径是否存在且为目录
if [ -z "$REPO_PATH" ]; then
    echo "Error: Unknown repository name '$REPO_NAME'."
    echo "Available repositories: ${!repo_map[@]}"
    exit 1
fi


# Change to the repository directory
cd "$REPO_PATH"
check_status "Failed to change to the repository directory"

# Ensure we're in a git repository
git rev-parse --is-inside-work-tree > /dev/null 2>&1
check_status "The provided path is not a git repository"


# 设置 trap，确保在脚本退出时切回原始分支
original_branch=$(git rev-parse --abbrev-ref HEAD)
trap 'log "切回初始分支 $original_branch..."; git checkout $original_branch; log "当前分支: $(git rev-parse --abbrev-ref HEAD)"' EXIT


check_status "Failed to get current branch name"
log "请确认当前分支：$original_branch"
read -p "分支是否正确，并继续执行脚本？(y/n): " confirm
if [[ $confirm != [yY] ]]; then
    log "脚本执行已取消。"
    exit 0
fi




# 舍弃本地所有修改，并拉取远程最新代码
echo "Fetching latest changes..."
handle_local_changes
git pull 

# 将当前分支合并到 develop 并推送 develop
echo "Merging $original_branch into develop..."
git checkout develop
check_status "Failed to checkout develop branch"
git merge $original_branch
check_status "Failed to merge $original_branch into develop"
git push origin develop 
check_status "Failed to push commit to remote develop branch"


# 将 develop 合并到 release 并推送 release
echo "Merging develop into release..."
git checkout release
check_status "Failed to checkout release branch"
git merge develop
check_status "Failed to merge develop into release"
git push origin release
check_status "Failed to push commit to remote release branch"



echo "Successfully merged branches and pushed changes."