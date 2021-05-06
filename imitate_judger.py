# -*- coding: utf-8 -*-
import sys
import os
import torch
import torchvision
from torchvision import transforms, datasets
import numpy as np
from PIL import Image


"""检查参数"""
if len(sys.argv) != 3:
    print("需要两个参数 python imitate_judger.py test_base.image test_user.image")
    sys.exit()

if not os.path.exists(sys.argv[1]):
    print("未找到第一个参数文件")
    sys.exit()

if not os.path.exists(sys.argv[2]):
    print("未找到第二个参数文件")
    sys.exit()

"""预定义超参数"""
IMG_SIZE = 256
INPUT_SIZE = 224
BATCH_SIZE = 256
EPOCHS_SIZE = 32
BASE_LR = 0.01
CLASSFIERS = 2
CUDA = torch.cuda.is_available()
DEVICE = torch.device('cuda' if CUDA else 'cpu')

"""预训练变换"""
preprocess_transform = transforms.Compose([
    transforms.Resize(IMG_SIZE),
    # transforms.CenterCrop(224),
    transforms.ToTensor(),
    transforms.Normalize([0.485, 0.456, 0.406], [0.229, 0.224, 0.225])
])

"""直接送入模型"""
image_test_PIL = Image.open(sys.argv[1]).convert('RGB')
image_user_PIL = Image.open(sys.argv[2]).convert('RGB')

"""与变换"""
image_test_tensor = preprocess_transform(image_test_PIL)
image_test_tensor.unsqueeze_(0)
image_test_tensor = image_test_tensor.to(DEVICE)
image_user_tensor = preprocess_transform(image_user_PIL)
image_user_tensor.unsqueeze_(0)
image_user_tensor = image_user_tensor.to(DEVICE)

"""定义模型"""
model = torchvision.models.resnet152(pretrained=True, num_classes=1000)
model.eval()

test_out = model(image_test_tensor)
user_out = model(image_user_tensor)

test_out_no_grad = test_out.detach()
user_out_no_grad = user_out.detach()

# test_out_no_grad[0] = test_out_no_grad[0] / np.linalg.norm(test_out_no_grad[0])
# user_out_no_grad[0] = user_out_no_grad[0] / np.linalg.norm(user_out_no_grad[0])
# dists = np.linalg.norm(test_out_no_grad[0] - user_out_no_grad[0])
# print(dists)

norm_diff = np.linalg.norm(test_out_no_grad[0] - user_out_no_grad[0])
norm1 = np.linalg.norm(test_out_no_grad[0])
norm2 = np.linalg.norm(user_out_no_grad[0])
print(1. - norm_diff / (norm1 + norm2))
